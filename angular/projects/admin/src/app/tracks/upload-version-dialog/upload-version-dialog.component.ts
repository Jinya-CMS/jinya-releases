import { Component, EventEmitter, Input, Output } from '@angular/core';
import { Application, Track, VersionService } from 'api';
import { FormControl, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'app-upload-version-dialog',
  templateUrl: './upload-version-dialog.component.html'
})
export class UploadVersionDialogComponent {
  @Input() application!: Application;
  @Input() track!: Track;
  isOpen = false;

  @Output() saved = new EventEmitter<void>();

  uploadVersionForm = new FormGroup({
    number: new FormControl<string>('', { nonNullable: true, validators: [Validators.required] }),
    artifact: new FormControl<File | null>(null, { nonNullable: true })
  });

  hasErrors = false;

  constructor(private versionService: VersionService) {}

  open(version = '') {
    this.uploadVersionForm.patchValue({ number: version });
    this.isOpen = true;
  }

  updateFile($event: Event) {
    // @ts-expect-error The event cannot be null
    const file = ($event.target as HTMLInputElement).files[0];
    this.uploadVersionForm.get('artifact')?.patchValue(file);
  }

  upload() {
    if (this.uploadVersionForm.valid) {
      const self = this;

      this.versionService
        .createNewVersion({
          body: this.uploadVersionForm.value.artifact!,
          versionNumber: this.uploadVersionForm.value.number!,
          trackId: this.track.id,
          applicationId: this.application.id
        })
        .subscribe({
          next() {
            self.saved.emit();
            self.isOpen = false;
          },
          error(e) {
            self.hasErrors = true;
            console.error(e);
          }
        });
    }
  }
}
