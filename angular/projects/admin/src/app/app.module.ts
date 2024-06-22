import { NgModule } from '@angular/core';
import { AppComponent } from './app.component';
import { NgIf } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterModule } from '@angular/router';
import { BrowserModule } from '@angular/platform-browser';
import { ApiModule } from '../../../api/src/lib/api.module';
import { provideHttpClient } from '@angular/common/http';
import { ApplicationsModule } from './applications/applications.module';

@NgModule({
  declarations: [AppComponent],
  imports: [
    ApiModule.forRoot({ rootUrl: '' }),
    BrowserModule,
    NgIf,
    ApplicationsModule,
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
  providers: [provideHttpClient()]
})
export class AppModule {}
