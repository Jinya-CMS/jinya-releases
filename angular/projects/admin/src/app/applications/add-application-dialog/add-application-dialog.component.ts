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
    })
  });

  hasErrors: boolean = false;

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

    // eslint-disable-next-line @typescript-eslint/no-this-alias
    const self = this;
    this.applicationService.createApplication({ body }).subscribe({
      next(value) {
        self.saved.emit(value);
        self.open = false;
        self.createApplicationForm.reset();
      },
      error(error) {
        self.hasErrors = true;
      }
    });
  }
}
