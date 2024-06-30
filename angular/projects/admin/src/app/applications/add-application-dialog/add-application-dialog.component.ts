import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Application, ApplicationService } from 'api';

@Component({
  selector: 'app-add-application-dialog',
  templateUrl: './add-application-dialog.component.html'
})
export class AddApplicationDialogComponent {
  createApplicationForm = new FormGroup({
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

  @Input() open = false;

  @Output() saved = new EventEmitter<Application>();

  constructor(private applicationService: ApplicationService) {}

  createApplication() {
    if (this.createApplicationForm.invalid) {
      return;
    }

    const body = {
      name: this.createApplicationForm.value.name!,
      slug: this.createApplicationForm.value.slug!,
      homepageTemplate: '<html></html>',
      trackpageTemplate: '<html></html>'
    };

    const self = this;
    this.applicationService.createApplication({ body }).subscribe({
      next(value) {
        if (self.createApplicationForm.value.logo) {
          const logoFile = self.createApplicationForm.value.logo;
          let response = null;
          switch (logoFile.type.replace('image/', '')) {
            case 'apng':
              response = self.applicationService.uploadApplicationLogo$Apng({ body: logoFile, id: value.id });
              break;
            case 'avif':
              response = self.applicationService.uploadApplicationLogo$Avif({ body: logoFile, id: value.id });
              break;
            case 'gif':
              response = self.applicationService.uploadApplicationLogo$Gif({ body: logoFile, id: value.id });
              break;
            case 'jpeg':
              response = self.applicationService.uploadApplicationLogo$Jpeg({ body: logoFile, id: value.id });
              break;
            case 'png':
              response = self.applicationService.uploadApplicationLogo$Png({ body: logoFile, id: value.id });
              break;
            case 'svg+xml':
              response = self.applicationService.uploadApplicationLogo$Xml({ body: logoFile, id: value.id });
              break;
            case 'webp':
              response = self.applicationService.uploadApplicationLogo$Webp({ body: logoFile, id: value.id });
              break;
          }
          if (response) {
            response.subscribe(() => {
              self.saved.emit(value);
              self.open = false;
              self.createApplicationForm.reset();
            });
          }
        } else {
          self.saved.emit(value);
          self.open = false;
          self.createApplicationForm.reset();
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
    this.createApplicationForm.get('logo')?.patchValue(file);
  }
}
