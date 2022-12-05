import {Component, OnInit} from '@angular/core';
import {MovieService} from "../_services/movie.service";
import {Movies} from "../_models/movies";
import {Movie} from "../_models/movie";
import {TokenStorageService} from "../_services/token-storage.service";

@Component({
  selector: 'app-home',
  templateUrl: './movies.component.html',
  styleUrls: ['./movies.component.css']
})
export class MoviesComponent implements OnInit {
  content?: string;
  movies: Movie[] = [];
  isAuthenticated: boolean = false


  constructor(private tokenService: TokenStorageService, private movieService: MovieService) {
  }

  ngOnInit(): void {
    this.isAuthenticated = false;
    if (this.tokenService.getToken() != null) {
      this.isAuthenticated = true;
    }
    let listFunc = this.movieService.listPublic()
    // if user is authenticated call the private API.
    if (this.isAuthenticated) {
      listFunc = this.movieService.list()
    }
    listFunc.subscribe({
      next: (data: Movies) => {
        this.movies = data.movies;
      },
      error: (err: { error: any; }) => {
        // @ts-ignore
        this.movies = JSON.parse(err.error).message;
      }
    });
  }

  makeAction(movie: Movie, action: string) {
    this.movieService.makeAction(movie.id, action).subscribe({
      next: () => {
        if (action == "like") {
          movie.likes++;
          movie.user_liked = true;
        } else {
          movie.hates++;
          movie.user_hated = true;
        }
      },
      error: (err: { error: any; }) => {
        // @ts-ignore
        this.movies = JSON.parse(err.error).message;
      }
    })
  }

  removeAction(movie: Movie, action: string) {
    this.movieService.removeAction(movie.id, action).subscribe({
      next: () => {
        if (action == "like") {
          movie.likes--;
          movie.user_liked = false;
        } else {
          movie.hates--;
          movie.user_hated = false;
        }
      },
      error: (err: { error: any; }) => {
        // @ts-ignore
        this.movies = JSON.parse(err.error).message;
      }
    })
  }

  sortMovies(sortType: string) {
    // if user is authenticated call the private API.
    let listFunc = this.movieService.listPublic(sortType)
    if (this.isAuthenticated) {
      listFunc = this.movieService.list(sortType)
    }
    listFunc.subscribe({
      next: (data: Movies) => {
        this.movies = data.movies;
      },
      error: (err: { error: any; }) => {
        this.movies = JSON.parse(err.error).message;
      }
    });
  }
}
