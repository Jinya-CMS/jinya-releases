import { Component } from '@angular/core';
import { Router, RouterLinkActive } from '@angular/router';
import { Location } from '@angular/common';
import { AuthenticationService } from '../authentication/authentication.service';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss',
  standalone: false
})
export class AppComponent {
  protected readonly RouterLinkActive = RouterLinkActive;

  constructor(
    protected authService: AuthenticationService,
    protected router: Router,
    protected location: Location
  ) {}
}
