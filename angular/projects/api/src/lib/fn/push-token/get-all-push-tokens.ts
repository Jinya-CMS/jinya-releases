/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { PushToken } from '../../models/push-token';

export interface GetAllPushTokens$Params {
}

export function getAllPushTokens(http: HttpClient, rootUrl: string, params?: GetAllPushTokens$Params, context?: HttpContext): Observable<StrictHttpResponse<Array<PushToken>>> {
  const rb = new RequestBuilder(rootUrl, getAllPushTokens.PATH, 'get');
  if (params) {
  }

  return http.request(
    rb.build({ responseType: 'json', accept: 'application/json', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return r as StrictHttpResponse<Array<PushToken>>;
    })
  );
}

getAllPushTokens.PATH = '/api/admin/push-token';
