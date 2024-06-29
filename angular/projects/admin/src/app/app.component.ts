import { Component } from '@angular/core';
import { Router, RouterLinkActive } from '@angular/router';
import { Location } from '@angular/common';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrl: './app.component.scss'
})
export class AppComponent {
  protected readonly RouterLinkActive = RouterLinkActive;

  constructor(
    protected router: Router,
    protected location: Location
  ) {}
}
