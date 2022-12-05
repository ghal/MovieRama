import {Component} from '@angular/core';
import {Movie} from "../_models/movie";
import {MovieService} from "../_services/movie.service";
import {Movies} from "../_models/movies";
import {ActivatedRoute} from "@angular/router";
import {TokenStorageService} from "../_services/token-storage.service";

@Component({
  selector: 'app-list-user-movies',
  templateUrl: './list-user-movies.component.html',
  styleUrls: ['./list-user-movies.component.scss']
})
export class ListUserMoviesComponent {
  constructor(private tokenService: TokenStorageService, private movieService: MovieService, private route: ActivatedRoute) {
  }

  movies: Movie[] = [];
  user_id: number = 0;
  isAuthenticated: boolean = false


  ngOnInit(): void {
    this.isAuthenticated = false;
    if (this.tokenService.getToken() != null) {
      this.isAuthenticated = true;
    }

    // @ts-ignore
    this.user_id = +this.route.snapshot.paramMap.get('user_id');

    let listFunc = this.movieService.listUserMoviesPublic(this.user_id)
    // if user is authenticated call the private API.
    if (this.isAuthenticated) {
      listFunc = this.movieService.listUserMovies(this.user_id)
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
    let listFunc = this.movieService.listUserMoviesPublic(this.user_id, sortType)
    // if user is authenticated call the private API.
    if (this.isAuthenticated) {
      listFunc = this.movieService.listUserMovies(this.user_id, sortType)
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
