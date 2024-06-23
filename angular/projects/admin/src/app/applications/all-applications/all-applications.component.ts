import { Component, Input, OnInit } from '@angular/core';
import { Application, ApplicationService, Track } from 'api';
import { Router } from '@angular/router';
import { ConfirmComponent } from '../../../ui/confirm/confirm.component';

enum ActiveTab {
  Details,
  Tracks
}

@Component({
  selector: 'app-all-applications',
  templateUrl: './all-applications.component.html'
})
export class AllApplicationsComponent implements OnInit {
  @Input() id?: string;

  applications: Application[] = [];
  selectedApplication: Application | null = null;
  activeTab = ActiveTab.Details;
  tracks: Track[] = [];
  loading = true;

  constructor(
    protected applicationService: ApplicationService,
    protected router: Router
  ) {}

  ngOnInit(): void {
    this.applicationService.getAllApplications().subscribe((value) => {
      this.applications = value;
      if (this.applications.length > 0) {
        if (!this.id) {
          this.router.navigateByUrl(`/application/${value[0].id}`);
        } else {
          this.selectedApplication = this.applications.find((app) => app.id === this.id) ?? this.applications[0];
          this.loading = false;
        }
      } else {
        this.loading = false;
      }
    });
  }

  appCreated(app: Application) {
    this.applications.push(app);
    this.selectedApplication = app;
  }

  appUpdated(app: Application) {
    const idx = this.applications.findIndex((application) => application.id === app.id);
    this.applications[idx].name = app.name;
    this.applications[idx].slug = app.slug;
    this.selectedApplication = app;
  }

  selectApp(app: Application) {
    this.selectedApplication = app;
  }

  deleteApp(deleteApplication: ConfirmComponent) {
    this.applicationService.deleteApplication({ id: this.selectedApplication?.id ?? '' }).subscribe(() => {
      deleteApplication.open = false;
      this.applications = this.applications.filter((application) => application.id !== this.selectedApplication?.id);
      if (this.applications.length > 0) {
        this.selectedApplication = this.applications[0];
      } else {
        this.selectedApplication = null;
      }
    });
  }

  protected readonly ActiveTab = ActiveTab;
  protected readonly location = location;
}
