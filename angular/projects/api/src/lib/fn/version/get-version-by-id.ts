/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Version } from '../../models/version';

export interface GetVersionById$Params {
  id: string;
  applicationId: string;
  trackId: string;
}

export function getVersionById(http: HttpClient, rootUrl: string, params: GetVersionById$Params, context?: HttpContext): Observable<StrictHttpResponse<Version>> {
  const rb = new RequestBuilder(rootUrl, getVersionById.PATH, 'get');
  if (params) {
    rb.path('id', params.id, {});
    rb.path('applicationId', params.applicationId, {});
    rb.path('trackId', params.trackId, {});
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

getVersionById.PATH = '/api/admin/application/{applicationId}/track/{trackId}/version/{id}';
