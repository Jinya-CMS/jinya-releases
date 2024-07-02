/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';


export interface CreateNewVersion$Params {
  versionNumber: string;
  applicationId: string;
  trackId: string;
      body: Blob
}

export function createNewVersion(http: HttpClient, rootUrl: string, params: CreateNewVersion$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
  const rb = new RequestBuilder(rootUrl, createNewVersion.PATH, 'post');
  if (params) {
    rb.path('versionNumber', params.versionNumber, {});
    rb.path('applicationId', params.applicationId, {});
    rb.path('trackId', params.trackId, {});
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

createNewVersion.PATH = '/api/admin/application/{applicationId}/track/{trackId}/version/{versionNumber}';
