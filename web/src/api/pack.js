
import axios from 'axios';

const http = axios.create();

export function ListPackTask(params){
    return http.get("/api/packTask", { params  });
}

export function CreatePackTask(params){
    return http.get("/api/packTask", { params });
}