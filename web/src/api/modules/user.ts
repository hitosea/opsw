import http from "../index";
import {User} from "../interface/user";

export const getUserInfo = () => {
    return http.get<User.Info>('user/info')
};

export const userLogin = (params: User.LoginReq) => {
    return http.post<User.Info>('user/login', params)
};

export const userReg = (params: User.RegReq) => {
    return http.post<User.Info>('user/reg', params)
};

export const getUserShareOptions = () => {
    return http.get<User.Options>('user/share-options')
};

export const userShareServer = (params:any) => {
    return http.post<User.ShareServerReq>('user/share-server', params)
};
