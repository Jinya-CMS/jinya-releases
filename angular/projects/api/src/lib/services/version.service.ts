/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { BaseService } from '../base-service';
import { ApiConfiguration } from '../api-configuration';
import { StrictHttpResponse } from '../strict-http-response';

import { createNewVersion } from '../fn/version/create-new-version';
import { CreateNewVersion$Params } from '../fn/version/create-new-version';
import { deleteVersion } from '../fn/version/delete-version';
import { DeleteVersion$Params } from '../fn/version/delete-version';
import { getAllVersions } from '../fn/version/get-all-versions';
import { GetAllVersions$Params } from '../fn/version/get-all-versions';
import { getVersionById } from '../fn/version/get-version-by-id';
import { GetVersionById$Params } from '../fn/version/get-version-by-id';
import { pushNewVersion } from '../fn/version/push-new-version';
import { PushNewVersion$Params } from '../fn/version/push-new-version';
import { Version } from '../models/version';

@Injectable({ providedIn: 'root' })
export class VersionService extends BaseService {
  constructor(config: ApiConfiguration, http: HttpClient) {
    super(config, http);
  }

  /** Path part for operation `getAllVersions()` */
  static readonly GetAllVersionsPath = '/api/admin/application/{applicationId}/track/{trackId}/version';

  /**
   * Get all versions.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getAllVersions()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllVersions$Response(params: GetAllVersions$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<Version>>> {
    return getAllVersions(this.http, this.rootUrl, params, context);
  }

  /**
   * Get all versions.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getAllVersions$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllVersions(params: GetAllVersions$Params, context?: HttpContext): Observable<Array<Version>> {
    return this.getAllVersions$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<Version>>): Array<Version> => r.body)
    );
  }

  /** Path part for operation `createNewVersion()` */
  static readonly CreateNewVersionPath = '/api/admin/application/{applicationId}/track/{trackId}/version/{versionNumber}';

  /**
   * Upload new version to the given application and track, if the version exists the binary is replaced.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `createNewVersion()` instead.
   *
   * This method sends `application/octet-stream` and handles request body of type `application/octet-stream`.
   */
  createNewVersion$Response(params: CreateNewVersion$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return createNewVersion(this.http, this.rootUrl, params, context);
  }

  /**
   * Upload new version to the given application and track, if the version exists the binary is replaced.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `createNewVersion$Response()` instead.
   *
   * This method sends `application/octet-stream` and handles request body of type `application/octet-stream`.
   */
  createNewVersion(params: CreateNewVersion$Params, context?: HttpContext): Observable<void> {
    return this.createNewVersion$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `getVersionById()` */
  static readonly GetVersionByIdPath = '/api/admin/application/{applicationId}/track/{trackId}/version/{id}';

  /**
   * Get version by id.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getVersionById()` instead.
   *
   * This method doesn't expect any request body.
   */
  getVersionById$Response(params: GetVersionById$Params, context?: HttpContext): Observable<StrictHttpResponse<Version>> {
    return getVersionById(this.http, this.rootUrl, params, context);
  }

  /**
   * Get version by id.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getVersionById$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getVersionById(params: GetVersionById$Params, context?: HttpContext): Observable<Version> {
    return this.getVersionById$Response(params, context).pipe(
      map((r: StrictHttpResponse<Version>): Version => r.body)
    );
  }

  /** Path part for operation `deleteVersion()` */
  static readonly DeleteVersionPath = '/api/admin/application/{applicationId}/track/{trackId}/version/{id}';

  /**
   * Delete version.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `deleteVersion()` instead.
   *
   * This method doesn't expect any request body.
   */
  deleteVersion$Response(params: DeleteVersion$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return deleteVersion(this.http, this.rootUrl, params, context);
  }

  /**
   * Delete version.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `deleteVersion$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  deleteVersion(params: DeleteVersion$Params, context?: HttpContext): Observable<void> {
    return this.deleteVersion$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `pushNewVersion()` */
  static readonly PushNewVersionPath = '/api/push/{applicationSlug}/{trackSlug}/{versionNumber}';

  /**
   * Push new version to the given application and track, if the version exists the binary is replaced.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `pushNewVersion()` instead.
   *
   * This method sends `application/octet-stream` and handles request body of type `application/octet-stream`.
   */
  pushNewVersion$Response(params: PushNewVersion$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return pushNewVersion(this.http, this.rootUrl, params, context);
  }

  /**
   * Push new version to the given application and track, if the version exists the binary is replaced.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `pushNewVersion$Response()` instead.
   *
   * This method sends `application/octet-stream` and handles request body of type `application/octet-stream`.
   */
  pushNewVersion(params: PushNewVersion$Params, context?: HttpContext): Observable<void> {
    return this.pushNewVersion$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

}
