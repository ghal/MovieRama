import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';

import {RegisterComponent} from './register/register.component';
import {LoginComponent} from './login/login.component';
import {MoviesComponent} from './movies/movies.component';
import {ListUserMoviesComponent} from "./list-user-movies/list-user-movies.component";
import {AddMovieComponent} from "./add-movie/add-movie.component";

// @ts-ignore
const routes: Routes = [
  {path: '', component: MoviesComponent},
  {path: 'users/:user_id/movies', component: ListUserMoviesComponent, pathMatch: 'full'},
  {path: 'login', component: LoginComponent},
  {path: 'register', component: RegisterComponent},
  {path: 'movies/add', component: AddMovieComponent},
  // {path: '', redirectTo: 'home', pathMatch: 'full'}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
