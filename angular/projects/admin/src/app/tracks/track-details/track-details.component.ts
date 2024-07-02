import { Component, Input, OnInit } from '@angular/core';
import { Application, ApplicationService, Track, TrackService, Version, VersionService } from 'api';
import { ConfirmComponent } from '../../../ui/confirm/confirm.component';
import { Router } from '@angular/router';
import { zip } from 'rxjs';
import { UploadVersionDialogComponent } from '../upload-version-dialog/upload-version-dialog.component';

@Component({
  selector: 'app-track-details',
  templateUrl: './track-details.component.html'
})
export class TrackDetailsComponent implements OnInit {
  versions!: Version[];
  application!: Application;
  track!: Track;
  loading = true;
  selectedVersion?: Version | null;

  @Input() trackId!: string;
  @Input() applicationId!: string;

  protected readonly location = location;

  constructor(
    private applicationService: ApplicationService,
    private trackService: TrackService,
    private versionService: VersionService,
    private router: Router
  ) {}

  ngOnInit() {
    zip(
      this.applicationService.getApplicationById({ id: this.applicationId }),
      this.trackService.getTrackById({
        id: this.trackId,
        applicationId: this.applicationId
      })
    ).subscribe(([application, track]) => {
      this.application = application;
      this.track = track;
      this.loadVersions();
    });
  }

  loadVersions() {
    this.versionService
      .getAllVersions({
        applicationId: this.applicationId,
        trackId: this.trackId
      })
      .subscribe((versions) => {
        this.versions = versions;
        this.loading = false;
      });
  }

  uploadDateToString(uploadDate: string) {
    const date = new Date(Date.parse(uploadDate));

    return date.toLocaleDateString();
  }

  openDeleteVersion(version: Version, deleteVersionConfirm: ConfirmComponent) {
    this.selectedVersion = version;
    deleteVersionConfirm.open = true;
  }

  deleteVersion(deleteVersionConfirm: ConfirmComponent) {
    this.versionService
      .deleteVersion({
        trackId: this.track.id,
        applicationId: this.application.id,
        id: this.selectedVersion!.id!
      })
      .subscribe(() => {
        deleteVersionConfirm.open = false;
        this.selectedVersion = null;
      });
  }

  openEditVersion(version: Version, uploadVersion: UploadVersionDialogComponent) {
    this.selectedVersion = version;
    uploadVersion.open(version.version);
  }
}
