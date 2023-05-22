import http from "../index";
import {User} from "../interface/user";

export const getUserInfo = () => {
    return http.get<User.Info>('user/info')
};

export const userLogin = (params: User.LoginReq) => {
    return http.post<User.Info>('user/login')
};

export const userReg = (params: User.RegReq) => {
    return http.post<User.Info>('user/reg')
};
