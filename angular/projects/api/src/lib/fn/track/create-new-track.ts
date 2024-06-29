/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { CreateTrack } from '../../models/create-track';
import { Track } from '../../models/track';

export interface CreateNewTrack$Params {
  applicationId: string;
      body?: CreateTrack
}

export function createNewTrack(http: HttpClient, rootUrl: string, params: CreateNewTrack$Params, context?: HttpContext): Observable<StrictHttpResponse<Track>> {
  const rb = new RequestBuilder(rootUrl, createNewTrack.PATH, 'post');
  if (params) {
    rb.path('applicationId', params.applicationId, {});
    rb.body(params.body, 'application/json');
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

createNewTrack.PATH = '/api/admin/application/{applicationId}/track';
