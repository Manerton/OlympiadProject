<?php

namespace App\Repositories;

use App\Components\ApiHelper;
use App\Services\ApiService;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Cookie;

class ApplicationRepository
{
    private ApiService $apiService;
    public function __construct(ApiService $apiService){
        $this->apiService = $apiService;
    }
    public function getByApiAll($page = 1, $limit = 10)
    {
        $token = request()->header('Authorization');
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response = $this->apiService->get(
            ApiHelper::APPLICATION_URL_API,
            [
                'page' => $page,
                'limit' => $limit
            ],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'] ? $response['data']['data'] : [];
    }
    public function getByEventId($id)
    {
        $token = request()->header('Authorization');
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response = $this->apiService->get(
            ApiHelper::APPLICATION_EVENT_URL_API . '/' . $id,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'] ? $response['data']['data'] : [];
    }
    public function getByApiId($id)
    {
        $token = request()->header('Authorization');
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response = $this->apiService->get(
            ApiHelper::APPLICATION_URL_API . '/' . $id,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'];
    }
    public function getCount()
    {
        $token = request()->header('Authorization');
        $token = !is_null($token) ? $token : json_decode(Cookie::get('username'))->token;
        $response = $this->apiService->get(
            ApiHelper::APPLICATION_COUNT_URL_API,
            [],
            [
                'Authorization' => "Bearer ". $token
            ]
        );
        return $response['data']['data'];
    }
}
