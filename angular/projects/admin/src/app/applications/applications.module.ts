import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { AllApplicationsComponent } from './all-applications/all-applications.component';

@NgModule({
  declarations: [],
  imports: [
    CommonModule,
    RouterModule.forChild([
      {
        path: 'application',
        component: AllApplicationsComponent
      },
      {
        path: 'application/:id',
        component: AllApplicationsComponent
      }
    ])
  ]
})
export class ApplicationsModule {}
