import {Component} from '@angular/core';
import {MovieService} from "../_services/movie.service";
import {CreateMovie} from "../_models/create-movie";
import {Router} from "@angular/router";

@Component({
  selector: 'app-add-movie',
  templateUrl: './add-movie.component.html',
  styleUrls: ['./add-movie.component.scss']
})
export class AddMovieComponent {
  form: any = {
    title: null,
    description: null
  };
  errorMessage = '';

  constructor(private movieService: MovieService, private router: Router) {
  }

  ngOnInit(): void {
  }

  onSubmit(): void {
    const mv: CreateMovie = this.form;
    this.movieService.create(mv).subscribe({
      next: (data: any) => {
        this.router.navigateByUrl('')
      },
      error: (err: { error: { message: string; }; }) => {
        this.errorMessage = err.error.message;
      }
    });
  }
}
