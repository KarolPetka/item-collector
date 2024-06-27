import { Injectable } from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
@Injectable({
  providedIn: 'root'
})
export class Repository {
  readonly url;

  constructor(private http: HttpClient) {
    this.url = "http://127.0.0.1:8000";
  }

  private createAuthHeaders(token: string): HttpHeaders {
    return new HttpHeaders({
      'Authorization': `Bearer ${token}`
    });
  }

  get(uri: string, token: string){
    return this.http.get(`${this.url}/${uri}`, { headers: this.createAuthHeaders(token) });
  }

  post(uri: string, payload: Object, token: string){
    return this.http.post(`${this.url}/${uri}`, payload, { headers: this.createAuthHeaders(token) });
  }

  postJson(url: string, data: any, token: string) {
    const headers = this.createAuthHeaders(token).append('Content-Type', 'application/json');
    return this.http.post(`${this.url}/${url}`, data, { headers });
  }

  postWithoutHeaders(uri: string, payload: Object){
    return this.http.post(`${this.url}/${uri}`, payload);
  }

  put(uri: string, item: any, token: string){
    return this.http.put(`${this.url}/${uri}`, item, { headers: this.createAuthHeaders(token) });
  }


  delete(uri:string, token: string){
    return this.http.delete(`${this.url}/${uri}`, { headers: this.createAuthHeaders(token) });
  }
}
