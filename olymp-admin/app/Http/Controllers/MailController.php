<?php

namespace App\Http\Controllers;

use App\Http\Requests\MailRequest;
use App\Services\MailService;

class MailController extends Controller
{
    private MailService $mailService;
    public function __construct(
        MailService $mailService
    )
    {
        $this->mailService = $mailService;
    }
    public function index(){
        return view('mail.index');
    }
    public function send(MailRequest $request){
        $data = $request->validated();
        $this->mailService->sendMessage($data['email'], $data['message']);
        return redirect()->route('mail.index');
    }
}
