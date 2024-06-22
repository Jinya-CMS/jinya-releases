/* tslint:disable */
/* eslint-disable */
import { HttpClient, HttpContext, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { filter, map } from 'rxjs/operators';
import { StrictHttpResponse } from '../../strict-http-response';
import { RequestBuilder } from '../../request-builder';


export interface UploadApplicationLogo$Gif$Params {
  id: string;
      body: Blob
}

export function uploadApplicationLogo$Gif(http: HttpClient, rootUrl: string, params: UploadApplicationLogo$Gif$Params, context?: HttpContext): Observable<StrictHttpResponse<void>> {
  const rb = new RequestBuilder(rootUrl, uploadApplicationLogo$Gif.PATH, 'post');
  if (params) {
    rb.path('id', params.id, {});
    rb.body(params.body, 'image/gif');
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

uploadApplicationLogo$Gif.PATH = '/api/admin/application/{id}/logo';
