<?php

namespace App\Http\Controllers;

use App\Components\ApiHelper;
use App\Http\Requests\LoginRequest;
use App\Services\ApiService;
use Illuminate\Support\Facades\Cookie;

class SiteController extends Controller
{
    private ApiService $apiService;
    public function __construct(
        ApiService $apiService
    )
    {
        $this->apiService = $apiService;
    }

    public function index()
    {
        return view('site/home');
    }
    public function login()
    {
        return view('site/login');
    }
    public function auth(LoginRequest $request){
        $data = $request->validated();
        $response = $this->apiService->post(ApiHelper::AUTH_URL_API, $data);
        if ($response['success']) {
            $token = $response['data']['data']['access_token'];
            Cookie::queue('username', json_encode([
                'email' => $data['email'],
                'token' => $token
            ]), 60 * 24 * 365);
            return redirect('/');
        }
        else {
            return view('site/login');
        }
    }
    public function logout(){
        Cookie::queue(Cookie::forget('username'));
        return redirect('/');
    }
}
