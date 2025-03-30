import Alpine from '/static/lib/alpine.js';
import { get, httpDelete, post } from '../../../lib/jinya-http.js';
import confirm from '../../lib/ui/confirm.js';
import { localize } from '../../../lib/jinya-alpine-tools.js';

Alpine.data('trackData', () => ({
  loading: true,
  versions: [],
  track: null,
  application: null,
  appId: null,
  trackId: null,
  uploadVersionData: {
    open: false,
    hasError: false,
    artifact: null,
    number: '',
    reset() {
      this.artifact = null;
      this.number = '';
    },
    updateFile(evt) {
      this.artifact = evt.target.files[0];
    },
  },
  async init() {
    const query = new URLSearchParams(window.location.search);
    this.appId = query.get('app');
    this.trackId = query.get('track');

    await Promise.all([
      (async () => {
        this.application = await get(`/api/admin/application/${this.appId}`);
      })(),
      (async () => {
        this.track = await get(`/api/admin/application/${this.appId}/track/${this.trackId}`);
      })(),
      (async () => {
        this.versions = await get(`/api/admin/application/${this.appId}/track/${this.trackId}/version`);
      })(),
    ]);

    this.loading = false;
  },
  toDateString(date) {
    if (date instanceof Date) {
      return date.toLocaleDateString();
    }

    if (typeof date === 'string') {
      return new Date(Date.parse(date)).toLocaleDateString();
    }

    if (typeof date === 'number') {
      return new Date(date).toLocaleDateString();
    }

    return date;
  },
  openUploadVersion() {
    this.uploadVersionData.reset();
    this.uploadVersionData.open = true;
  },
  async deleteVersion(version) {
    if (
      await confirm({
        title: localize({ key: 'delete-version-title' }),
        message: localize({
          key: 'delete-version-message',
          values: { versionNumber: version.version, appName: this.application.name, trackName: this.track.name },
        }),
        declineLabel: localize({ key: 'delete-version-decline' }),
        approveLabel: localize({ key: 'delete-version-confirm' }),
      })
    ) {
      try {
        await httpDelete(`/api/admin/application/${this.appId}/track/${this.trackId}/version/${version.id}`);
        this.versions = this.versions.filter((v) => v.id !== version.id);
      } catch (error) {
        console.error(error);
      }
    }
  },
  async uploadVersion() {
    try {
      await post(
        `/api/admin/application/${this.appId}/track/${this.trackId}/version/${this.uploadVersionData.number}`,
        this.uploadVersionData.artifact,
      );
      this.uploadVersionData.open = false;
      this.versions = await get(`/api/admin/application/${this.appId}/track/${this.trackId}/version`);
    } catch (error) {
      console.error(error);
      this.uploadVersionData.hasError = true;
    }
  },
}));
