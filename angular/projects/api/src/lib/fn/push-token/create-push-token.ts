/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';

import { CreateToken } from '../../models/create-token';
import { PushToken } from '../../models/push-token';

export interface CreatePushToken$Params {
      body?: CreateToken
}

export function createPushToken(http: HttpClient, rootUrl: string, params?: CreatePushToken$Params, context?: HttpContext): Observable<StrictHttpResponse<PushToken>> {
  const rb = new RequestBuilder(rootUrl, createPushToken.PATH, 'post');
  if (params) {
    rb.body(params.body, 'application/json');
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

createPushToken.PATH = '/api/admin/push-token';
