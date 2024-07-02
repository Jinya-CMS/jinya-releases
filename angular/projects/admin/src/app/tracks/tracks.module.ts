import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { TrackDetailsComponent } from './track-details/track-details.component';
import { RouterModule } from '@angular/router';
import { ConfirmComponent } from '../../ui/confirm/confirm.component';
import { UiModule } from '../../ui/ui.module';
import { UploadVersionDialogComponent } from './upload-version-dialog/upload-version-dialog.component';
import { ReactiveFormsModule } from '@angular/forms';

@NgModule({
  declarations: [TrackDetailsComponent, UploadVersionDialogComponent],
  imports: [
    CommonModule,
    RouterModule.forChild([
      {
        path: 'application/:applicationId/track/:trackId',
        component: TrackDetailsComponent
      }
    ]),
    ConfirmComponent,
    UiModule,
    ReactiveFormsModule
  ]
})
export class TracksModule {}
