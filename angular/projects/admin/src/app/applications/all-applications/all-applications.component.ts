import { Component, Input, OnInit } from '@angular/core';
import { Application, ApplicationService, Track, TrackService, VersionService } from 'api';
import { Router } from '@angular/router';
import { ConfirmComponent } from '../../../ui/confirm/confirm.component';
import { EditTrackDialogComponent } from '../edit-track-dialog/edit-track-dialog.component';

enum ActiveTab {
  Details,
  Tracks,
  PushTokens
}

@Component({
  selector: 'app-all-applications',
  templateUrl: './all-applications.component.html',
  styleUrl: 'all-applications.component.scss'
})
export class AllApplicationsComponent implements OnInit {
  @Input() id?: string;

  applications: Application[] = [];
  selectedApplication: Application | null = null;
  selectedTrack: Track | null = null;
  activeTab = ActiveTab.Details;
  tracks: Track[] = [];
  loading = true;
  trackHasVersions: Record<string, boolean> = {};
  logoError = false;
  timestamp = new Date().getTime();
  newToken: string | null = null;
  tokenCopied = false;

  constructor(
    protected applicationService: ApplicationService,
    protected trackService: TrackService,
    protected versionService: VersionService,
    protected router: Router
  ) {}

  ngOnInit(): void {
    this.applicationService.getAllApplications().subscribe((value) => {
      this.applications = value;
      if (this.applications.length > 0) {
        if (!this.id) {
          this.router.navigateByUrl(`/application/${value[0].id}`);
        } else {
          this.selectApp(this.applications.find((app) => app.id === this.id) ?? this.applications[0]);
          this.loading = false;
        }
      } else {
        this.loading = false;
      }
    });
  }

  appCreated(app: Application) {
    this.applications.push(app);
    this.selectApp(app);
  }

  appUpdated(app: Application) {
    const idx = this.applications.findIndex((application) => application.id === app.id);
    this.applications[idx].name = app.name;
    this.applications[idx].slug = app.slug;
    this.selectedApplication = app;
    this.timestamp = new Date().getTime();
  }

  selectApp(app: Application) {
    this.logoError = false;
    this.timestamp = new Date().getTime();
    this.selectedApplication = app;
    this.trackService
      .getAllTracks({
        applicationId: app.id
      })
      .subscribe((tracks) => {
        this.tracks = tracks;
        for (const track of tracks) {
          this.versionService
            .getAllVersions({
              applicationId: this.selectedApplication!.id,
              trackId: track.id
            })
            .subscribe((versions) => (this.trackHasVersions[track.id] = versions.length > 0));
        }
      });
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

  trackAdded(track: Track) {
    this.tracks.push(track);
  }

  updateTrack(track: Track) {
    if (track.isDefault) {
      this.tracks = this.tracks.map((t) => ({ ...t, isDefault: false }));
    }

    const idx = this.tracks.findIndex((t) => track.id === t.id);
    this.tracks[idx].name = track.name;
    this.tracks[idx].slug = track.slug;
    this.tracks[idx].isDefault = track.isDefault;
  }

  editTrack(track: Track, editTrackDialog: EditTrackDialogComponent) {
    editTrackDialog.open(this.selectedApplication!, track);
  }

  deleteTrack(deleteTrackConfirm: ConfirmComponent) {
    this.trackService
      .deleteTrack({
        applicationId: this.selectedApplication!.id,
        id: this.selectedTrack!.id
      })
      .subscribe(() => {
        deleteTrackConfirm.open = false;
        this.tracks = this.tracks.filter((t) => t.id !== this.selectedTrack!.id);
        this.selectedTrack = null;
      });
  }

  openDeleteTrack(track: Track, deleteTrackConfirm: ConfirmComponent) {
    this.selectedTrack = track;
    deleteTrackConfirm.open = true;
  }

  protected readonly ActiveTab = ActiveTab;
  protected readonly location = location;

  createPushToken(createPushTokenDialog: ConfirmComponent) {
    this.tokenCopied = false;
    this.applicationService.createToken({ id: this.selectedApplication!.id }).subscribe((value) => {
      this.newToken = value.token;
      createPushTokenDialog.open = false;
    });
  }

  resetTokens(resetPushTokenDialog: ConfirmComponent) {
    this.applicationService.resetTokens({ id: this.selectedApplication!.id }).subscribe(() => {
      resetPushTokenDialog.open = false;
    });
  }

  async copyToken() {
    await navigator.clipboard.writeText(this.newToken ?? '');
    this.tokenCopied = true;
  }
}
