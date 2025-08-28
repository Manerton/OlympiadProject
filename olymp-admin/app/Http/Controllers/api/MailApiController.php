<?php

namespace App\Http\Controllers\api;

use App\Http\Controllers\Controller;
use App\Http\Requests\MailRequest;
use App\Services\MailService;

class MailApiController extends Controller
{
    private MailService $mailService;
    public function __construct(
        MailService $mailService
    )
    {
        $this->mailService = $mailService;
    }
    public function index(){
        return response()->json([]);
    }
    public function send(MailRequest $request){
        $data = $request->validated();
        return response()->json([$this->mailService->sendMessage($data['email'], $data['message'])]);
    }
}
