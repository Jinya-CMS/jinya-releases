/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { PushToken } from '../../models/push-token';

export interface GetPushTokenById$Params {
  id: string;
}

export function getPushTokenById(http: HttpClient, rootUrl: string, params: GetPushTokenById$Params, context?: HttpContext): Observable<StrictHttpResponse<PushToken>> {
  const rb = new RequestBuilder(rootUrl, getPushTokenById.PATH, 'get');
  if (params) {
    rb.path('id', params.id, {});
  }

  return http.request(
    rb.build({ responseType: 'json', accept: 'application/json', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<PushToken>;
    })
  );
}

getPushTokenById.PATH = '/api/admin/push-token/{id}';
