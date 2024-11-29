import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Application, ApplicationService } from 'api';

@Component({
  selector: 'app-edit-application-dialog',
  templateUrl: './edit-application-dialog.component.html',
  standalone: false
})
export class EditApplicationDialogComponent {
  editApplicationForm = new FormGroup({
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
    logo: new FormControl<File | null>(null, { nonNullable: true })
  });

  hasErrors = false;

  isOpen = false;

  @Input() selectedApplication!: Application;

  @Output() saved = new EventEmitter<Application>();

  constructor(private applicationService: ApplicationService) {}

  open() {
    this.isOpen = true;
    this.editApplicationForm.reset({
      name: this.selectedApplication.name,
      slug: this.selectedApplication.slug
    });
    this.hasErrors = false;
  }

  editApplication() {
    if (this.editApplicationForm.invalid) {
      return;
    }

    const body = {
      ...this.selectedApplication,
      name: this.editApplicationForm.value.name!,
      slug: this.editApplicationForm.value.slug!
    };

    const self = this;
    this.applicationService.updateApplication({ id: this.selectedApplication.id, body }).subscribe({
      next() {
        if (self.editApplicationForm.value.logo) {
          const logoFile = self.editApplicationForm.value.logo;
          let response = null;
          switch (logoFile.type.replace('image/', '')) {
            case 'apng':
              response = self.applicationService.uploadApplicationLogo$Apng({
                body: logoFile,
                id: self.selectedApplication.id
              });
              break;
            case 'avif':
              response = self.applicationService.uploadApplicationLogo$Avif({
                body: logoFile,
                id: self.selectedApplication.id
              });
              break;
            case 'gif':
              response = self.applicationService.uploadApplicationLogo$Gif({
                body: logoFile,
                id: self.selectedApplication.id
              });
              break;
            case 'jpeg':
              response = self.applicationService.uploadApplicationLogo$Jpeg({
                body: logoFile,
                id: self.selectedApplication.id
              });
              break;
            case 'png':
              response = self.applicationService.uploadApplicationLogo$Png({
                body: logoFile,
                id: self.selectedApplication.id
              });
              break;
            case 'svg+xml':
              response = self.applicationService.uploadApplicationLogo$Xml({
                body: logoFile,
                id: self.selectedApplication.id
              });
              break;
            case 'webp':
              response = self.applicationService.uploadApplicationLogo$Webp({
                body: logoFile,
                id: self.selectedApplication.id
              });
              break;
          }
          if (response) {
            response.subscribe(() => {
              self.saved.emit({
                ...self.selectedApplication,
                name: self.editApplicationForm.value.name!,
                slug: self.editApplicationForm.value.slug!
              });
              self.isOpen = false;
            });
          }
        } else {
          self.saved.emit({
            ...self.selectedApplication,
            name: self.editApplicationForm.value.name!,
            slug: self.editApplicationForm.value.slug!
          });
          self.isOpen = false;
        }
      },
      error() {
        self.hasErrors = true;
      }
    });
  }

  updateFile($event: Event) {
    // @ts-expect-error The event cannot be null
    const file = ($event.target as HTMLInputElement).files[0];
    this.editApplicationForm.get('logo')?.patchValue(file);
  }
}
