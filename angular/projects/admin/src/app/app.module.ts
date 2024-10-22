import { NgModule, provideExperimentalZonelessChangeDetection } from '@angular/core';
import { AppComponent } from './app.component';
import { NgIf } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterModule } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { ApiModule } from 'api';
import { provideHttpClient } from '@angular/common/http';
import { ApplicationsModule } from './applications/applications.module';
import { UiModule } from '../ui/ui.module';
import { OAuthModule } from 'angular-oauth2-oidc';
import { AuthenticationModule } from '../authentication/authentication.module';
import { TracksModule } from './tracks/tracks.module';

@NgModule({
  declarations: [AppComponent],
  imports: [
    ApiModule.forRoot({ rootUrl: '' }),
    BrowserModule,
    NgIf,
    ApplicationsModule,
    TracksModule,
    UiModule,
    OAuthModule.forRoot(),
    AuthenticationModule.forRoot(),
    RouterLink,
    RouterLinkActive,
    RouterModule.forRoot(
      [
        {
          path: '**',
          redirectTo: 'application'
        }
      ],
      {
        bindToComponentInputs: true
      }
    )
  ],
  bootstrap: [AppComponent],
  providers: [provideExperimentalZonelessChangeDetection(), provideHttpClient()]
})
export class AppModule {}
