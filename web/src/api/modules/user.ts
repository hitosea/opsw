import http from "../index";
import {User} from "../interface/user";

export const getUserInfo = () => {
    return http.get<User.Info>('user/info')
};
