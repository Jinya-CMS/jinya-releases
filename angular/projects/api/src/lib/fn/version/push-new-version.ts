/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';


export interface PushNewVersion$Params {
  versionNumber: string;
  applicationSlug: string;
  trackSlug: string;
      body: Blob
}

export function pushNewVersion(http: HttpClient, rootUrl: string, params: PushNewVersion$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
  const rb = new RequestBuilder(rootUrl, pushNewVersion.PATH, 'post');
  if (params) {
    rb.path('versionNumber', params.versionNumber, {});
    rb.path('applicationSlug', params.applicationSlug, {});
    rb.path('trackSlug', params.trackSlug, {});
    rb.body(params.body, 'application/octet-stream');
  }

  return http.request(
    rb.build({ responseType: 'text', accept: '*/*', context })
  ).pipe(
    filter((r: any): r is HttpResponse<any> => r instanceof HttpResponse),
    map((r: HttpResponse<any>) => {
      return (r as HttpResponse<any>).clone({ body: undefined }) as StrictHttpResponse<void>;
    })
  );
}

pushNewVersion.PATH = '/api/push/{applicationSlug}/{trackSlug}/{versionNumber}';
