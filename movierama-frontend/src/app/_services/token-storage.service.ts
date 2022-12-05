import {Injectable} from '@angular/core';

const TOKEN_KEY = 'auth-token';
const USER_KEY = 'auth-user';

@Injectable({
  providedIn: 'root'
})
export class TokenStorageService {
  constructor() {
  }

  signOut(): void {
    // @ts-ignore
    window.sessionStorage.clear();
  }

  public saveToken(token: string): void {
    // @ts-ignore
    window.sessionStorage.removeItem(TOKEN_KEY);
    // @ts-ignore
    window.sessionStorage.setItem(TOKEN_KEY, token);
  }

  public getToken(): string | null {
    // @ts-ignore
    return window.sessionStorage.getItem(TOKEN_KEY);
  }

  public saveUser(user: any): void {
    // @ts-ignore
    window.sessionStorage.removeItem(USER_KEY);
    // @ts-ignore
    window.sessionStorage.setItem(USER_KEY, JSON.stringify(user));
  }

  public getUser(): any {
    // @ts-ignore
    const user = window.sessionStorage.getItem(USER_KEY);
    if (user) {
      // @ts-ignore
      return JSON.parse(user);
    }

    return {};
  }
}
