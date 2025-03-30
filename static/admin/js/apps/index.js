import Alpine from '/static/lib/alpine.js';
import { get, httpDelete, post, put } from '../../../lib/jinya-http.js';
import confirm from '../../lib/ui/confirm.js';
import { localize } from '../../../lib/jinya-alpine-tools.js';

Alpine.data('appsData', () => ({
  loading: true,
  selectedApp: null,
  apps: [],
  tracks: [],
  createPushTokenOpen: false,
  createTrackOpen: false,
  editTrackOpen: false,
  trackHasVersions: new Set(),
  activeTab: 'details',
  createApp: {
    open: false,
    hasError: false,
    slug: '',
    name: '',
    reset() {
      this.slug = '';
      this.name = '';
      this.hasError = false;
    },
  },
  editApp: {
    open: false,
    hasError: false,
    slug: '',
    name: '',
    reset(app) {
      this.slug = app.slug;
      this.name = app.name;
      this.hasError = false;
    },
  },
  createTrackData: {
    open: false,
    hasError: false,
    slug: '',
    name: '',
    isDefault: false,
    reset() {
      this.slug = '';
      this.name = '';
      this.isDefault = false;
      this.hasError = false;
    },
  },
  editTrackData: {
    open: false,
    hasError: false,
    id: '',
    slug: '',
    name: '',
    isDefault: false,
    reset(track) {
      this.id = track.id;
      this.slug = track.slug;
      this.name = track.name;
      this.isDefault = track.isDefault;
      this.hasError = false;
    },
  },
  displayPushToken: {
    open: false,
    token: '',
    copied: false,
  },
  get activeLink() {
    return `${location.origin}/${this.selectedApp.slug}`;
  },
  async init() {
    this.apps = await get('/api/admin/application');
    await this.selectApp(this.apps[0]);
    this.loading = false;
  },
  async selectApp(app) {
    this.selectedApp = app;
    this.tracks = await get(`/api/admin/application/${this.selectedApp.id}/track`);
    const versionsByTrack = new Set();
    for await (const track of this.tracks) {
      if ((await get(`/api/admin/application/${this.selectedApp.id}/track/${track.id}/version`)).length > 0) {
        versionsByTrack.add(track.id);
      }
    }
    this.trackHasVersions = versionsByTrack;
  },
  getTrackLink(track) {
    return `${location.origin}/${this.selectedApp.slug}/${track.slug}`;
  },
  async deleteApplication() {
    if (
      await confirm({
        title: localize({ key: 'delete-app-title' }),
        message: localize({ key: 'delete-app-message', values: this.selectedApp }),
        declineLabel: localize({ key: 'delete-app-decline' }),
        approveLabel: localize({ key: 'delete-app-confirm' }),
        negative: true,
      })
    ) {
      try {
        await httpDelete(`/api/admin/application/${this.selectedApp.id}`);
        this.apps = this.apps.filter((app) => app.id !== this.selectedApp.id);
        await this.selectApp(this.apps[0]);
      } catch (error) {
        console.error(error);
      }
    }
  },
  openCreateApplication() {
    this.createApp.reset();
    this.createApp.open = true;
  },
  openEditApplication() {
    this.editApp.reset(this.selectedApp);
    this.editApp.open = true;
  },
  openCreateTrack() {
    this.createTrackData.open = true;
    this.createTrackData.reset();
  },
  openEditTrack(track) {
    this.editTrackData.open = true;
    this.selectedTrack = track;
    this.editTrackData.reset(track);
  },
  async deleteTrack(track) {
    if (
      await confirm({
        title: localize({ key: 'delete-track-title' }),
        message: localize({
          key: 'delete-track-message',
          values: {
            trackName: track.name,
            appName: this.selectedApp.name,
          },
        }),
        declineLabel: localize({ key: 'delete-track-decline' }),
        approveLabel: localize({ key: 'delete-track-confirm' }),
        negative: true,
      })
    ) {
      try {
        await httpDelete(`/api/admin/application/${this.selectedApp.id}/track/${track.id}`);
        this.tracks = this.tracks.filter((t) => t.id !== track.id);
      } catch (error) {
        console.error(error);
      }
    }
  },
  async resetPushTokens() {
    if (
      await confirm({
        title: localize({ key: 'reset-push-tokens-title' }),
        message: localize({ key: 'reset-push-tokens-message', values: this.selectedApp }),
        declineLabel: localize({ key: 'reset-push-tokens-decline' }),
        approveLabel: localize({ key: 'reset-push-tokens-confirm' }),
        negative: true,
      })
    ) {
      try {
        await httpDelete(`/api/admin/application/${this.selectedApp.id}/token`);
      } catch (error) {
        console.error(error);
      }
    }
  },
  async createPushToken() {
    if (
      await confirm({
        title: localize({ key: 'create-push-token-title' }),
        message: localize({ key: 'create-push-token-message', values: this.selectedApp }),
        declineLabel: localize({ key: 'create-push-token-decline' }),
        approveLabel: localize({ key: 'create-push-token-confirm' }),
      })
    ) {
      const token = await post(`/api/admin/application/${this.selectedApp.id}/token`);
      this.displayPushToken.open = true;
      this.displayPushToken.token = token.token;
      this.displayPushToken.copied = false;
    }
  },
  async createApplication() {
    try {
      const newApp = await post('/api/admin/application', this.createApp);
      this.apps.push(newApp);
      await this.selectApp(newApp);
      this.createApp.hasError = false;
      this.createApp.open = false;
    } catch (e) {
      console.error('Error creating application', e);
      this.createApp.hasError = true;
    }
  },
  async editApplication() {
    try {
      await put(`/api/admin/application/${this.selectedApp.id}`, this.editApp);
      const selectedIdx = this.apps.findIndex((app) => app.id === this.selectedApp.id);
      this.apps[selectedIdx].name = this.editApp.name;
      this.apps[selectedIdx].slug = this.editApp.slug;
      this.selectedApp.name = this.editApp.name;
      this.selectedApp.slug = this.editApp.slug;
      this.editApp.hasError = false;
      this.editApp.open = false;
    } catch (e) {
      console.error('Error creating application', e);
      this.createApp.hasError = true;
    }
  },
  async createTrack() {
    try {
      const newTrack = await post(`/api/admin/application/${this.selectedApp.id}/track`, this.createTrackData);
      if (this.createTrackData.isDefault) {
        this.tracks.forEach((track) => {
          track.isDefault = false;
        });
      }
      this.tracks.push(newTrack);
      this.createTrackData.hasError = false;
      this.createTrackData.open = false;
    } catch (e) {
      console.error('Error creating track', e);
      this.createTrackData.hasError = true;
    }
  },
  async editTrack() {
    try {
      await put(`/api/admin/application/${this.selectedApp.id}/track/${this.editTrackData.id}`, this.editTrackData);
      const selectedIdx = this.tracks.findIndex((app) => app.id === this.editTrackData.id);
      if (this.editTrackData.isDefault) {
        this.tracks.forEach((track) => {
          track.isDefault = false;
        });
      }
      this.tracks[selectedIdx].name = this.editTrackData.name;
      this.tracks[selectedIdx].slug = this.editTrackData.slug;
      this.tracks[selectedIdx].isDefault = this.editTrackData.isDefault;
      this.editTrackData.hasError = false;
      this.editTrackData.open = false;
    } catch (e) {
      console.error('Error creating application', e);
      this.createApp.hasError = true;
    }
  },
  async copyToken() {
    await navigator.clipboard.writeText(this.displayPushToken.token ?? '');
    this.displayPushToken.copied = true;
  },
}));
