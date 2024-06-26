import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Application, Track, TrackService } from 'api';

@Component({
  selector: 'app-edit-track-dialog',
  templateUrl: './edit-track-dialog.component.html'
})
export class EditTrackDialogComponent {
  editTrackForm = new FormGroup({
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

  hasErrors: boolean = false;

  @Input() isOpen = false;
  @Input() selectedApplication!: Application;
  track!: Track;

  @Output() saved = new EventEmitter<Track>();

  constructor(private trackService: TrackService) {}

  updateTrack() {
    if (this.editTrackForm.invalid) {
      return;
    }

    // eslint-disable-next-line @typescript-eslint/no-this-alias
    const self = this;

    this.trackService
      .updateTrack({
        id: this.track.id,
        applicationId: this.selectedApplication.id,
        body: {
          name: this.editTrackForm.value.name ?? '',
          slug: this.editTrackForm.value.slug ?? '',
          isDefault: this.editTrackForm.value.isDefault ?? false
        }
      })
      .subscribe({
        next() {
          self.saved.emit({
            ...self.track,
            name: self.editTrackForm.value.name!,
            slug: self.editTrackForm.value.slug!,
            isDefault: self.editTrackForm.value.isDefault!
          });
          self.isOpen = false;
        },
        error() {
          self.hasErrors = true;
        }
      });
  }

  open(application: Application, track: Track) {
    this.track = track;
    this.selectedApplication = application;
    this.isOpen = true;
    this.editTrackForm.reset({
      name: this.track.name,
      slug: this.track.slug,
      isDefault: this.track.isDefault
    });
    this.hasErrors = false;
  }
}
