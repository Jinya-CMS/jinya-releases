/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';


export interface UploadApplicationLogo$Avif$Params {
  id: string;
      body: Blob
}

export function uploadApplicationLogo$Avif(http: HttpClient, rootUrl: string, params: UploadApplicationLogo$Avif$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
  const rb = new RequestBuilder(rootUrl, uploadApplicationLogo$Avif.PATH, 'post');
  if (params) {
    rb.path('id', params.id, {});
    rb.body(params.body, 'image/avif');
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

uploadApplicationLogo$Avif.PATH = '/api/admin/application/{id}/logo';
