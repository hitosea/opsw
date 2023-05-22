import {defineStore} from 'pinia';
import {WsMsg, WsState} from './interface';
import {UserStore} from "./user";
import utils from "../utils/utils";
import {CONST} from "./constant";

export const WsStore = defineStore({
    id: 'WsState',
    state: (): WsState => ({
        uid: 0,
        rid: null,
        ws: null,
        msg: {},
        timeout: null,
        random: "",
        openNum: 0,
        listener: {},
        watch: false
    }),
    actions: {
        connection() {
            const userStore = UserStore()
            //
            clearTimeout(this.timeout);
            if (this.ws) {
                this.ws.close();
                this.ws = null;
            }
            if (userStore.info.id === 0) {
                return;
            }
            //
            let url = window.location.origin + "/ws";
            url = url.replace("https://", "wss://");
            url = url.replace("http://", "ws://");
            url += `?token=${userStore.info.token}`;
            //
            const wgLog = true;
            const random = utils.randomString(16);
            this.random = random;
            //
            this.ws = new WebSocket(url);
            this.ws.onopen = async (e) => {
                wgLog && console.log("[WS] Open", e, utils.formatDate())
                this.openNum++;
            };
            this.ws.onclose = async (e) => {
                wgLog && console.log("[WS] Close", e, utils.formatDate())
                this.ws = null;
                //
                clearTimeout(this.timeout);
                this.timeout = setTimeout(() => {
                    random === this.random && this.connection();
                }, 3000);
            };
            this.ws.onerror = async (e) => {
                wgLog && console.log("[WS] Error", e, utils.formatDate())
                this.ws = null;
                //
                clearTimeout(this.timeout);
                this.timeout = setTimeout(() => {
                    random === this.random && this.connection();
                }, 3000);
            };
            this.ws.onmessage = async (e) => {
                wgLog && console.log("[WS] Message", e);
                const wsMsg: WsMsg = utils.jsonParse(e.data);
                this.msg = wsMsg;
                //
                if (wsMsg.action === CONST.WsOnline) {
                    // 连接成功
                    if (wsMsg.data.own === 1) {
                        this.rid = wsMsg.data.rid
                    }
                }
                Object.values(this.listener).forEach(call => {
                    if (typeof call === "function") {
                        try {
                            call(wsMsg);
                        } catch (err) {
                            wgLog && console.log("[WS] Callerr", err);
                        }
                    }
                });
            }
        },

        listener(name, callback) {
            if (typeof callback === "function") {
                this.listener[name] = callback;
            } else {
                this.listener[name] && delete this.listener[name];
            }
        },

        close() {
            if (this.ws) {
                this.ws.close();
                this.ws = null;
            }
        }
    },
});
