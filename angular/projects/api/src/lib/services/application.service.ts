/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { BaseService } from '../base-service';
import { ApiConfiguration } from '../api-configuration';
import { StrictHttpResponse } from '../strict-http-response';

import { Application } from '../models/application';
import { createApplication } from '../fn/application/create-application';
import { CreateApplication$Params } from '../fn/application/create-application';
import { createToken } from '../fn/application/create-token';
import { CreateToken$Params } from '../fn/application/create-token';
import { deleteApplication } from '../fn/application/delete-application';
import { DeleteApplication$Params } from '../fn/application/delete-application';
import { getAllApplications } from '../fn/application/get-all-applications';
import { GetAllApplications$Params } from '../fn/application/get-all-applications';
import { getApplicationById } from '../fn/application/get-application-by-id';
import { GetApplicationById$Params } from '../fn/application/get-application-by-id';
import { PushToken } from '../models/push-token';
import { resetTokens } from '../fn/application/reset-tokens';
import { ResetTokens$Params } from '../fn/application/reset-tokens';
import { updateApplication } from '../fn/application/update-application';
import { UpdateApplication$Params } from '../fn/application/update-application';
import { uploadApplicationLogo$Apng } from '../fn/application/upload-application-logo-apng';
import { UploadApplicationLogo$Apng$Params } from '../fn/application/upload-application-logo-apng';
import { uploadApplicationLogo$Avif } from '../fn/application/upload-application-logo-avif';
import { UploadApplicationLogo$Avif$Params } from '../fn/application/upload-application-logo-avif';
import { uploadApplicationLogo$Gif } from '../fn/application/upload-application-logo-gif';
import { UploadApplicationLogo$Gif$Params } from '../fn/application/upload-application-logo-gif';
import { uploadApplicationLogo$Jpeg } from '../fn/application/upload-application-logo-jpeg';
import { UploadApplicationLogo$Jpeg$Params } from '../fn/application/upload-application-logo-jpeg';
import { uploadApplicationLogo$Png } from '../fn/application/upload-application-logo-png';
import { UploadApplicationLogo$Png$Params } from '../fn/application/upload-application-logo-png';
import { uploadApplicationLogo$Webp } from '../fn/application/upload-application-logo-webp';
import { UploadApplicationLogo$Webp$Params } from '../fn/application/upload-application-logo-webp';
import { uploadApplicationLogo$Xml } from '../fn/application/upload-application-logo-xml';
import { UploadApplicationLogo$Xml$Params } from '../fn/application/upload-application-logo-xml';

@Injectable({ providedIn: 'root' })
export class ApplicationService extends BaseService {
  constructor(config: ApiConfiguration, http: HttpClient) {
    super(config, http);
  }

  /** Path part for operation `getAllApplications()` */
  static readonly GetAllApplicationsPath = '/api/admin/application';

  /**
   * Get all applications.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getAllApplications()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllApplications$Response(params?: GetAllApplications$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<Application>>> {
    return getAllApplications(this.http, this.rootUrl, params, context);
  }

  /**
   * Get all applications.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getAllApplications$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllApplications(params?: GetAllApplications$Params, context?: HttpContext): Observable<Array<Application>> {
    return this.getAllApplications$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<Application>>): Array<Application> => r.body)
    );
  }

  /** Path part for operation `createApplication()` */
  static readonly CreateApplicationPath = '/api/admin/application';

  /**
   * Create application.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `createApplication()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  createApplication$Response(params?: CreateApplication$Params, context?: HttpContext): Observable<StrictHttpResponse<Application>> {
    return createApplication(this.http, this.rootUrl, params, context);
  }

  /**
   * Create application.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `createApplication$Response()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  createApplication(params?: CreateApplication$Params, context?: HttpContext): Observable<Application> {
    return this.createApplication$Response(params, context).pipe(
      map((r: StrictHttpResponse<Application>): Application => r.body)
    );
  }

  /** Path part for operation `getApplicationById()` */
  static readonly GetApplicationByIdPath = '/api/admin/application/{id}';

  /**
   * Get application by id.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getApplicationById()` instead.
   *
   * This method doesn't expect any request body.
   */
  getApplicationById$Response(params: GetApplicationById$Params, context?: HttpContext): Observable<StrictHttpResponse<Application>> {
    return getApplicationById(this.http, this.rootUrl, params, context);
  }

  /**
   * Get application by id.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getApplicationById$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getApplicationById(params: GetApplicationById$Params, context?: HttpContext): Observable<Application> {
    return this.getApplicationById$Response(params, context).pipe(
      map((r: StrictHttpResponse<Application>): Application => r.body)
    );
  }

  /** Path part for operation `updateApplication()` */
  static readonly UpdateApplicationPath = '/api/admin/application/{id}';

  /**
   * Update application.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `updateApplication()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  updateApplication$Response(params: UpdateApplication$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return updateApplication(this.http, this.rootUrl, params, context);
  }

  /**
   * Update application.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `updateApplication$Response()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  updateApplication(params: UpdateApplication$Params, context?: HttpContext): Observable<void> {
    return this.updateApplication$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `deleteApplication()` */
  static readonly DeleteApplicationPath = '/api/admin/application/{id}';

  /**
   * Delete Application.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `deleteApplication()` instead.
   *
   * This method doesn't expect any request body.
   */
  deleteApplication$Response(params: DeleteApplication$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return deleteApplication(this.http, this.rootUrl, params, context);
  }

  /**
   * Delete Application.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `deleteApplication$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  deleteApplication(params: DeleteApplication$Params, context?: HttpContext): Observable<void> {
    return this.deleteApplication$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `createToken()` */
  static readonly CreateTokenPath = '/api/admin/application/{id}/token';

  /**
   * Create token.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `createToken()` instead.
   *
   * This method doesn't expect any request body.
   */
  createToken$Response(params: CreateToken$Params, context?: HttpContext): Observable<StrictHttpResponse<PushToken>> {
    return createToken(this.http, this.rootUrl, params, context);
  }

  /**
   * Create token.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `createToken$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  createToken(params: CreateToken$Params, context?: HttpContext): Observable<PushToken> {
    return this.createToken$Response(params, context).pipe(
      map((r: StrictHttpResponse<PushToken>): PushToken => r.body)
    );
  }

  /** Path part for operation `resetTokens()` */
  static readonly ResetTokensPath = '/api/admin/application/{id}/token';

  /**
   * Resets all allowed push tokens for the given application.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `resetTokens()` instead.
   *
   * This method doesn't expect any request body.
   */
  resetTokens$Response(params: ResetTokens$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return resetTokens(this.http, this.rootUrl, params, context);
  }

  /**
   * Resets all allowed push tokens for the given application.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `resetTokens$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  resetTokens(params: ResetTokens$Params, context?: HttpContext): Observable<void> {
    return this.resetTokens$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `uploadApplicationLogo()` */
  static readonly UploadApplicationLogoPath = '/api/admin/application/{id}/logo';

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `uploadApplicationLogo$Apng()` instead.
   *
   * This method sends `image/apng` and handles request body of type `image/apng`.
   */
  uploadApplicationLogo$Apng$Response(params: UploadApplicationLogo$Apng$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return uploadApplicationLogo$Apng(this.http, this.rootUrl, params, context);
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `uploadApplicationLogo$Apng$Response()` instead.
   *
   * This method sends `image/apng` and handles request body of type `image/apng`.
   */
  uploadApplicationLogo$Apng(params: UploadApplicationLogo$Apng$Params, context?: HttpContext): Observable<void> {
    return this.uploadApplicationLogo$Apng$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `uploadApplicationLogo$Avif()` instead.
   *
   * This method sends `image/avif` and handles request body of type `image/avif`.
   */
  uploadApplicationLogo$Avif$Response(params: UploadApplicationLogo$Avif$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return uploadApplicationLogo$Avif(this.http, this.rootUrl, params, context);
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `uploadApplicationLogo$Avif$Response()` instead.
   *
   * This method sends `image/avif` and handles request body of type `image/avif`.
   */
  uploadApplicationLogo$Avif(params: UploadApplicationLogo$Avif$Params, context?: HttpContext): Observable<void> {
    return this.uploadApplicationLogo$Avif$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `uploadApplicationLogo$Gif()` instead.
   *
   * This method sends `image/gif` and handles request body of type `image/gif`.
   */
  uploadApplicationLogo$Gif$Response(params: UploadApplicationLogo$Gif$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return uploadApplicationLogo$Gif(this.http, this.rootUrl, params, context);
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `uploadApplicationLogo$Gif$Response()` instead.
   *
   * This method sends `image/gif` and handles request body of type `image/gif`.
   */
  uploadApplicationLogo$Gif(params: UploadApplicationLogo$Gif$Params, context?: HttpContext): Observable<void> {
    return this.uploadApplicationLogo$Gif$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `uploadApplicationLogo$Jpeg()` instead.
   *
   * This method sends `image/jpeg` and handles request body of type `image/jpeg`.
   */
  uploadApplicationLogo$Jpeg$Response(params: UploadApplicationLogo$Jpeg$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return uploadApplicationLogo$Jpeg(this.http, this.rootUrl, params, context);
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `uploadApplicationLogo$Jpeg$Response()` instead.
   *
   * This method sends `image/jpeg` and handles request body of type `image/jpeg`.
   */
  uploadApplicationLogo$Jpeg(params: UploadApplicationLogo$Jpeg$Params, context?: HttpContext): Observable<void> {
    return this.uploadApplicationLogo$Jpeg$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `uploadApplicationLogo$Png()` instead.
   *
   * This method sends `image/png` and handles request body of type `image/png`.
   */
  uploadApplicationLogo$Png$Response(params: UploadApplicationLogo$Png$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return uploadApplicationLogo$Png(this.http, this.rootUrl, params, context);
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `uploadApplicationLogo$Png$Response()` instead.
   *
   * This method sends `image/png` and handles request body of type `image/png`.
   */
  uploadApplicationLogo$Png(params: UploadApplicationLogo$Png$Params, context?: HttpContext): Observable<void> {
    return this.uploadApplicationLogo$Png$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `uploadApplicationLogo$Xml()` instead.
   *
   * This method sends `image/svg+xml` and handles request body of type `image/svg+xml`.
   */
  uploadApplicationLogo$Xml$Response(params: UploadApplicationLogo$Xml$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return uploadApplicationLogo$Xml(this.http, this.rootUrl, params, context);
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `uploadApplicationLogo$Xml$Response()` instead.
   *
   * This method sends `image/svg+xml` and handles request body of type `image/svg+xml`.
   */
  uploadApplicationLogo$Xml(params: UploadApplicationLogo$Xml$Params, context?: HttpContext): Observable<void> {
    return this.uploadApplicationLogo$Xml$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `uploadApplicationLogo$Webp()` instead.
   *
   * This method sends `image/webp` and handles request body of type `image/webp`.
   */
  uploadApplicationLogo$Webp$Response(params: UploadApplicationLogo$Webp$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return uploadApplicationLogo$Webp(this.http, this.rootUrl, params, context);
  }

  /**
   * Upload application logo.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `uploadApplicationLogo$Webp$Response()` instead.
   *
   * This method sends `image/webp` and handles request body of type `image/webp`.
   */
  uploadApplicationLogo$Webp(params: UploadApplicationLogo$Webp$Params, context?: HttpContext): Observable<void> {
    return this.uploadApplicationLogo$Webp$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

}
