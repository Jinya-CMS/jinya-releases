/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

import { BaseService } from '../base-service';
import { ApiConfiguration } from '../api-configuration';
import { StrictHttpResponse } from '../strict-http-response';

import { createPushToken } from '../fn/push-token/create-push-token';
import { CreatePushToken$Params } from '../fn/push-token/create-push-token';
import { deletePushToken } from '../fn/push-token/delete-push-token';
import { DeletePushToken$Params } from '../fn/push-token/delete-push-token';
import { getAllPushTokens } from '../fn/push-token/get-all-push-tokens';
import { GetAllPushTokens$Params } from '../fn/push-token/get-all-push-tokens';
import { getPushTokenById } from '../fn/push-token/get-push-token-by-id';
import { GetPushTokenById$Params } from '../fn/push-token/get-push-token-by-id';
import { PushToken } from '../models/push-token';
import { updatePushToken } from '../fn/push-token/update-push-token';
import { UpdatePushToken$Params } from '../fn/push-token/update-push-token';

@Injectable({ providedIn: 'root' })
export class PushTokenService extends BaseService {
  constructor(config: ApiConfiguration, http: HttpClient) {
    super(config, http);
  }

  /** Path part for operation `getAllPushTokens()` */
  static readonly GetAllPushTokensPath = '/api/admin/push-token';

  /**
   * Get all push tokens.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getAllPushTokens()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllPushTokens$Response(params?: GetAllPushTokens$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<PushToken>>> {
    return getAllPushTokens(this.http, this.rootUrl, params, context);
  }

  /**
   * Get all push tokens.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getAllPushTokens$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getAllPushTokens(params?: GetAllPushTokens$Params, context?: HttpContext): Observable<Array<PushToken>> {
    return this.getAllPushTokens$Response(params, context).pipe(
      map((r: StrictHttpResponse<Array<PushToken>>): Array<PushToken> => r.body)
    );
  }

  /** Path part for operation `createPushToken()` */
  static readonly CreatePushTokenPath = '/api/admin/push-token';

  /**
   * Create push token.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `createPushToken()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  createPushToken$Response(params?: CreatePushToken$Params, context?: HttpContext): Observable<StrictHttpResponse<PushToken>> {
    return createPushToken(this.http, this.rootUrl, params, context);
  }

  /**
   * Create push token.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `createPushToken$Response()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  createPushToken(params?: CreatePushToken$Params, context?: HttpContext): Observable<PushToken> {
    return this.createPushToken$Response(params, context).pipe(
      map((r: StrictHttpResponse<PushToken>): PushToken => r.body)
    );
  }

  /** Path part for operation `getPushTokenById()` */
  static readonly GetPushTokenByIdPath = '/api/admin/push-token/{id}';

  /**
   * Get push token by id.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `getPushTokenById()` instead.
   *
   * This method doesn't expect any request body.
   */
  getPushTokenById$Response(params: GetPushTokenById$Params, context?: HttpContext): Observable<StrictHttpResponse<PushToken>> {
    return getPushTokenById(this.http, this.rootUrl, params, context);
  }

  /**
   * Get push token by id.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `getPushTokenById$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  getPushTokenById(params: GetPushTokenById$Params, context?: HttpContext): Observable<PushToken> {
    return this.getPushTokenById$Response(params, context).pipe(
      map((r: StrictHttpResponse<PushToken>): PushToken => r.body)
    );
  }

  /** Path part for operation `updatePushToken()` */
  static readonly UpdatePushTokenPath = '/api/admin/push-token/{id}';

  /**
   * Update push token.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `updatePushToken()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  updatePushToken$Response(params: UpdatePushToken$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return updatePushToken(this.http, this.rootUrl, params, context);
  }

  /**
   * Update push token.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `updatePushToken$Response()` instead.
   *
   * This method sends `application/json` and handles request body of type `application/json`.
   */
  updatePushToken(params: UpdatePushToken$Params, context?: HttpContext): Observable<void> {
    return this.updatePushToken$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

  /** Path part for operation `deletePushToken()` */
  static readonly DeletePushTokenPath = '/api/admin/push-token/{id}';

  /**
   * Delete push token.
   *
   *
   *
   * This method provides access to the full `HttpResponse`, allowing access to response headers.
   * To access only the response body, use `deletePushToken()` instead.
   *
   * This method doesn't expect any request body.
   */
  deletePushToken$Response(params: DeletePushToken$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
    return deletePushToken(this.http, this.rootUrl, params, context);
  }

  /**
   * Delete push token.
   *
   *
   *
   * This method provides access only to the response body.
   * To access the full response (for headers, for example), `deletePushToken$Response()` instead.
   *
   * This method doesn't expect any request body.
   */
  deletePushToken(params: DeletePushToken$Params, context?: HttpContext): Observable<void> {
    return this.deletePushToken$Response(params, context).pipe(
      map((r: StrictHttpResponse<void>): void => r.body)
    );
  }

}
