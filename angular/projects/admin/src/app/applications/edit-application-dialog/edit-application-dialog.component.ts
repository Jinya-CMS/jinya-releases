import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Application, ApplicationService } from 'api';

@Component({
  selector: 'app-edit-application-dialog',
  templateUrl: './edit-application-dialog.component.html'
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
    })
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
        self.saved.emit({
          ...self.selectedApplication,
          name: self.editApplicationForm.value.name!,
          slug: self.editApplicationForm.value.slug!
        });
        self.isOpen = false;
      },
      error() {
        self.hasErrors = true;
      }
    });
  }
}
