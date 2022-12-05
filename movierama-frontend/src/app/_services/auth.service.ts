import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable} from 'rxjs';

const AUTH_API = 'http://localhost:1323/api/v1/auth/';

const httpOptions = {
  headers: new HttpHeaders({'Content-Type': 'application/json'})
};

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  constructor(private http: HttpClient) {
  }

  login(username: string, password: string): Observable<any> {
    return this.http.post(AUTH_API + 'login', {
      username: username,
      password: password
    }, httpOptions);
  }

  register(username: string, password: string, firstName: string, lastName: string): Observable<any> {
    return this.http.post(AUTH_API + 'register', {
      username: username,
      password: password,
      first_name: firstName,
      last_name: lastName
    }, httpOptions);
  }
}
