import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Application, Track, TrackService } from 'api';

@Component({
  selector: 'app-add-track-dialog',
  templateUrl: './add-track-dialog.component.html',
  standalone: false
})
export class AddTrackDialogComponent {
  createTrackForm = new FormGroup({
    name: new FormControl<string>('', { nonNullable: true, validators: [Validators.required] }),
    slug: new FormControl<string>('', {
      nonNullable: true,
      validators: [
        Validators.required,
        (control) => {
          if (/[a-zA-Z0-9_-]/gm.test(control.value)) {
            return null;
          }

          return {
            pattern: false
          };
        }
      ]
    }),
    isDefault: new FormControl<boolean>(false)
  });

  hasErrors = false;

  @Input() open = false;
  @Input() selectedApplication!: Application;

  @Output() saved = new EventEmitter<Track>();

  constructor(private trackService: TrackService) {}

  createTrack() {
    if (this.createTrackForm.invalid) {
      return;
    }

    const self = this;
    this.trackService
      .createNewTrack({
        applicationId: this.selectedApplication?.id,
        body: {
          name: this.createTrackForm.value.name ?? '',
          slug: this.createTrackForm.value.slug ?? '',
          isDefault: this.createTrackForm.value.isDefault ?? false
        }
      })
      .subscribe({
        next(track) {
          self.saved.emit(track);
          self.createTrackForm.reset();
          self.open = false;
        },
        error() {
          self.hasErrors = true;
        }
      });
  }
}
