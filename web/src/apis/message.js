import axios from 'axios';

export function apiListMsg(type, state) {
    return axios.get(`/community/message?type=${type}&state=${state}`);
}

export function apiClearUnReadMsg(type){
    return axios.post(`/community/message/UnReadMsg/${type}`);
}

