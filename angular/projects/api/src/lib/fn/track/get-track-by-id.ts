/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { Track } from '../../models/track';

export interface GetTrackById$Params {
  id: string;
  applicationId: string;
}

export function getTrackById(http: HttpClient, rootUrl: string, params: GetTrackById$Params, context?: HttpContext): Observable<StrictHttpResponse<Track>> {
  const rb = new RequestBuilder(rootUrl, getTrackById.PATH, 'get');
  if (params) {
    rb.path('id', params.id, {});
    rb.path('applicationId', params.applicationId, {});
  }

  return http.request(
    rb.build({ responseType: 'json', accept: 'application/json', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Track>;
    })
  );
}

getTrackById.PATH = '/api/admin/application/{applicationId}/track/{id}';
