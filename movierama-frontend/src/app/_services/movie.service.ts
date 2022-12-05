import {Injectable} from '@angular/core';
import {HttpClient, HttpHeaders} from "@angular/common/http";
import {Observable} from "rxjs";
import {Movies} from "../_models/movies";
import {CreateMovie} from "../_models/create-movie";

const PUBLIC_MOVIE_API = 'http://localhost:1323/';
const MOVIE_API = 'http://localhost:1323/api/v1/';

const httpOptions = {
  headers: new HttpHeaders({'Content-Type': 'application/json'})
};

@Injectable({
  providedIn: 'root'
})
export class MovieService {

  constructor(private http: HttpClient) {
  }

  listPublic(sortType: string = "date"): Observable<Movies> {
    return this.http.get<Movies>(PUBLIC_MOVIE_API + "movies?sort=" + sortType, httpOptions);
  }

  list(sortType: string = "date"): Observable<Movies> {
    return this.http.get<Movies>(MOVIE_API + "movies?sort=" + sortType, httpOptions);
  }

  create(movie: CreateMovie): Observable<any> {
    return this.http.post(MOVIE_API + "movies", movie);
  }

  listUserMovies(user_id: number, sortType: string = "date"): Observable<Movies> {
    return this.http.get<Movies>(MOVIE_API + "users/" + user_id + "/movies?sort=" + sortType, httpOptions);
  }

  listUserMoviesPublic(user_id: number, sortType: string = "date"): Observable<Movies> {
    return this.http.get<Movies>(PUBLIC_MOVIE_API + "users/" + user_id + "/movies?sort=" + sortType, httpOptions);
  }

  makeAction(movie_id: number, action: string): Observable<any> {
    return this.http.post(MOVIE_API + "movies/" + movie_id + "/action/" + action, httpOptions);
  }

  removeAction(movie_id: number, action: string): Observable<any> {
    return this.http.post(MOVIE_API + "movies/" + movie_id + "/remove_action/" + action, httpOptions);
  }
}
