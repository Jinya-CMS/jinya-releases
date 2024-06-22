/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { BaseService } from '../base-service';
import { ApiConfiguration } from '../api-configuration';
import { StrictHttpResponse } from '../strict-http-response';

import { createNewTrack } from '../fn/track/create-new-track';
import { CreateNewTrack$Params } from '../fn/track/create-new-track';
import { deleteTrack } from '../fn/track/delete-track';
import { DeleteTrack$Params } from '../fn/track/delete-track';
import { getAllTracks } from '../fn/track/get-all-tracks';
import { GetAllTracks$Params } from '../fn/track/get-all-tracks';
import { getTrackById } from '../fn/track/get-track-by-id';
import { GetTrackById$Params } from '../fn/track/get-track-by-id';
import { Track } from '../models/track';
import { updateTrack } from '../fn/track/update-track';
import { UpdateTrack$Params } from '../fn/track/update-track';

@Injectable({ providedIn: 'root' })
export class TrackService extends BaseService {
  constructor(config: ApiConfiguration, http: HttpClient) {
    super(config, http);
  }

  /** Path part for operation `getAllTracks()` */
  static readonly GetAllTracksPath = '/api/admin/application/{applicationId}/track';

  /**
   * Get all tracks.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getAllTracks()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllTracks$Response(params: GetAllTracks$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<Track>>> {
    return getAllTracks(this.http, this.rootUrl, params, context);
  }

  /**
   * Get all tracks.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getAllTracks$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllTracks(params: GetAllTracks$Params, context?: HttpContext): Observable<Array<Track>> {
    return this.getAllTracks$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<Track>>): Array<Track> => r.body)
    );
  }

  /** Path part for operation `createNewTrack()` */
  static readonly CreateNewTrackPath = '/api/admin/application/{applicationId}/track';

  /**
   * Create track.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `createNewTrack()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  createNewTrack$Response(params: CreateNewTrack$Params, context?: HttpContext): Observable<StrictHttpResponse<Track>> {
    return createNewTrack(this.http, this.rootUrl, params, context);
  }

  /**
   * Create track.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `createNewTrack$Response()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  createNewTrack(params: CreateNewTrack$Params, context?: HttpContext): Observable<Track> {
    return this.createNewTrack$Response(params, context).pipe(
      map((r: StrictHttpResponse<Track>): Track => r.body)
    );
  }

  /** Path part for operation `getTrackById()` */
  static readonly GetTrackByIdPath = '/api/admin/application/{applicationId}/track/{id}';

  /**
   * Get track by id.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getTrackById()` instead.
   *
   * This method doesn't expect any request body.
   */
  getTrackById$Response(params: GetTrackById$Params, context?: HttpContext): Observable<StrictHttpResponse<Track>> {
    return getTrackById(this.http, this.rootUrl, params, context);
  }

  /**
   * Get track by id.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getTrackById$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getTrackById(params: GetTrackById$Params, context?: HttpContext): Observable<Track> {
    return this.getTrackById$Response(params, context).pipe(
      map((r: StrictHttpResponse<Track>): Track => r.body)
    );
  }

  /** Path part for operation `updateTrack()` */
  static readonly UpdateTrackPath = '/api/admin/application/{applicationId}/track/{id}';

  /**
   * Update track.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `updateTrack()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  updateTrack$Response(params: UpdateTrack$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return updateTrack(this.http, this.rootUrl, params, context);
  }

  /**
   * Update track.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `updateTrack$Response()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  updateTrack(params: UpdateTrack$Params, context?: HttpContext): Observable<void> {
    return this.updateTrack$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `deleteTrack()` */
  static readonly DeleteTrackPath = '/api/admin/application/{applicationId}/track/{id}';

  /**
   * Delete track.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `deleteTrack()` instead.
   *
   * This method doesn't expect any request body.
   */
  deleteTrack$Response(params: DeleteTrack$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return deleteTrack(this.http, this.rootUrl, params, context);
  }

  /**
   * Delete track.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `deleteTrack$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  deleteTrack(params: DeleteTrack$Params, context?: HttpContext): Observable<void> {
    return this.deleteTrack$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

}
