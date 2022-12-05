import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ListUserMoviesComponent } from './list-user-movies.component';

describe('ListUserMoviesComponent', () => {
  let component: ListUserMoviesComponent;
  let fixture: ComponentFixture<ListUserMoviesComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ ListUserMoviesComponent ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ListUserMoviesComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
