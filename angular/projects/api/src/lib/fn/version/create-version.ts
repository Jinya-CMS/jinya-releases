/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { CreateVersion } from '../../models/create-version';
import { Version } from '../../models/version';

export interface CreateVersion$Params {
  applicationId: string;
  trackId: string;
      body?: CreateVersion
}

export function createVersion(http: HttpClient, rootUrl: string, params: CreateVersion$Params, context?: HttpContext): Observable<StrictHttpResponse<Version>> {
  const rb = new RequestBuilder(rootUrl, createVersion.PATH, 'post');
  if (params) {
    rb.path('applicationId', params.applicationId, {});
    rb.path('trackId', params.trackId, {});
    rb.body(params.body, 'application/json');
  }

  return http.request(
    rb.build({ responseType: 'json', accept: 'application/json', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Version>;
    })
  );
}

createVersion.PATH = '/api/admin/application/{applicationId}/track/{trackId}/version';
