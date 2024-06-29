/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';


export interface UploadVersionArtifact$Params {
  id: string;
  applicationId: string;
  trackId: string;
      body: Blob
}

export function uploadVersionArtifact(http: HttpClient, rootUrl: string, params: UploadVersionArtifact$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
  const rb = new RequestBuilder(rootUrl, uploadVersionArtifact.PATH, 'post');
  if (params) {
    rb.path('id', params.id, {});
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

uploadVersionArtifact.PATH = '/api/admin/application/{applicationId}/track/{trackId}/version/{id}/file';
