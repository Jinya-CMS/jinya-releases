import { NgModule } from '@angular/core';
import { CommonModule, NgOptimizedImage } from '@angular/common';
import { RouterLink, RouterModule } from '@angular/router';
import { AllApplicationsComponent } from './all-applications/all-applications.component';
import { AddApplicationDialogComponent } from './add-application-dialog/add-application-dialog.component';
import { ReactiveFormsModule } from '@angular/forms';
import { CircleCheck, CircleX, LucideAngularModule } from 'lucide-angular';
import { EditApplicationDialogComponent } from './edit-application-dialog/edit-application-dialog.component';
import { UiModule } from '../../ui/ui.module';
import { ConfirmComponent } from '../../ui/confirm/confirm.component';
import { AddTrackDialogComponent } from './add-track-dialog/add-track-dialog.component';
import { EditTrackDialogComponent } from './edit-track-dialog/edit-track-dialog.component';

@NgModule({
  declarations: [
    AddApplicationDialogComponent,
    AllApplicationsComponent,
    EditApplicationDialogComponent,
    AddTrackDialogComponent,
    EditTrackDialogComponent
  ],
  imports: [
    RouterLink,
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
    ]),
    ReactiveFormsModule,
    LucideAngularModule.pick({ CircleX, CircleCheck }),
    UiModule,
    ConfirmComponent,
    NgOptimizedImage
  ]
})
export class ApplicationsModule {}
