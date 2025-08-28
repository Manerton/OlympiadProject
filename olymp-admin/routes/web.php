<?php

use app\controllers\EmailController;
use App\Http\Controllers\ApplicationController;
use App\Http\Controllers\EventController;
use App\Http\Controllers\MailController;
use App\Http\Controllers\ParticipantController;
use App\Http\Controllers\ReportController;
use App\Http\Controllers\SchoolController;
use App\Http\Controllers\SiteController;
use App\Http\Controllers\UserController;
use Illuminate\Support\Facades\Route;
use Prometheus\CollectorRegistry;
use Prometheus\RenderTextFormat;
use Prometheus\Storage\InMemory;
use Prometheus\Storage\Redis;
/*
|--------------------------------------------------------------------------
| Web Routes
|--------------------------------------------------------------------------
|
| Here is where you can register web routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| contains the "web" middleware group. Now create something great!
|
*/


Route::get('/login', [SiteController::class, 'login'])->name('login');
Route::post('/login', [SiteController::class, 'auth'])->name('auth');
Route::post('/logout', [SiteController::class, 'logout'])->name('logout');
Route::group(['middleware' => 'auth.custom'], function() {
    Route::get('/', [SiteController::class, 'index'])->name('homepage');
    Route::get('/school/index/{page?}', [SchoolController::class, 'index'])->name('school.index');
    Route::get('/school/create', [SchoolController::class, 'create'])->name('school.create');
    Route::post('/school/store', [SchoolController::class, 'store'])->name('school.store');
    Route::get('/school/edit/{id}', [SchoolController::class, 'edit'])->name('school.edit');
    Route::post('/school/update/{id}', [SchoolController::class, 'update'])->name('school.update');
    Route::delete('/school/delete/{id}', [SchoolController::class, 'delete'])->name('school.delete');
    Route::get('/school/show/{id}', [SchoolController::class, 'show'])->name('school.show');

    Route::get('/user/index/{page?}', [UserController::class, 'index'])->name('user.index');
    Route::get('/user/create', [UserController::class, 'create'])->name('user.create');
    Route::post('/user/store', [UserController::class, 'store'])->name('user.store');
    Route::get('/user/edit/{id}', [UserController::class, 'edit'])->name('user.edit');
    Route::post('/user/update/{id}', [UserController::class, 'update'])->name('user.update');
    Route::delete('/user/delete/{id}', [UserController::class, 'delete'])->name('user.delete');
    Route::get('/user/show/{id}', [UserController::class, 'show'])->name('user.show');
    Route::post('/user/revoke/{id}', [UserController::class, 'revoke'])->name('user.revoke');

    Route::get('/participant/index/{page?}', [ParticipantController::class, 'index'])->name('participant.index');
    Route::get('/participant/create', [ParticipantController::class, 'create'])->name('participant.create');
    Route::post('/participant/store', [ParticipantController::class, 'store'])->name('participant.store');
    Route::get('/participant/edit/{id}', [ParticipantController::class, 'edit'])->name('participant.edit');
    Route::post('/participant/update/{id}', [ParticipantController::class, 'update'])->name('participant.update');
    Route::delete('/participant/delete/{id}', [ParticipantController::class, 'delete'])->name('participant.delete');
    Route::get('/participant/show/{id}', [ParticipantController::class, 'show'])->name('participant.show');

    Route::get('/application/index/{page?}', [ApplicationController::class, 'index'])->name('application.index');
    Route::get('/application/create', [ApplicationController::class, 'create'])->name('application.create');
    Route::post('/application/store', [ApplicationController::class, 'store'])->name('application.store');
    Route::get('/application/edit/{id}', [ApplicationController::class, 'edit'])->name('application.edit');
    Route::post('/application/update/{id}', [ApplicationController::class, 'update'])->name('application.update');
    Route::delete('/application/delete/{id}', [ApplicationController::class, 'delete'])->name('application.delete');
    Route::get('/application/show/{id}', [ApplicationController::class, 'show'])->name('application.show');
    Route::post('/application/confirm/{id}', [ApplicationController::class, 'confirm'])->name('application.confirm');
    Route::post('/application/reject/{id}', [ApplicationController::class, 'reject'])->name('application.reject');

    Route::get('/event/index/{page?}', [EventController::class, 'index'])->name('event.index');
    Route::get('/event/show/{id}', [EventController::class, 'show'])->name('event.show');
    Route::delete('/event/delete/{id}', [EventController::class, 'delete'])->name('event.delete');
    Route::get('/event/task/{id}', [EventController::class, 'task'])->name('event.task');
    Route::get('/event/attendance/{id}', [EventController::class, 'attendance'])->name('event.attendance');
    Route::get('/event/point/{id}', [EventController::class, 'point'])->name('event.point');
    Route::get('/event/synchronize/{id}', [EventController::class, 'synchronize'])->name('event.synchronize');
    Route::post('/event/add-task/{id}', [EventController::class, 'addTask'])->name('event.add-task');
    Route::post('/event/change-attendance', [EventController::class, 'changeAttendance'])->name('event.change-attendance');
    Route::delete('/event/delete-task/{id}', [EventController::class, 'deleteTask'])->name('event.delete-task');
    Route::post('/event/change-score', [EventController::class, 'changeScore'])->name('event.change-score');
    Route::get('/event/prize-score/{id}', [EventController::class, 'prizeScore'])->name('event.prize-score');
    Route::post('/event/set-prize-score/{id}', [EventController::class, 'setPrizeScore'])->name('event.set-prize-score');
    Route::get('/event/appeal/{id}', [EventController::class, 'appeal'])->name('event.appeal');

    Route::get('/report/index', [ReportController::class, 'index'])->name('report.index');
    Route::get('/report/download/{id}', [ReportController::class, 'download'])->name('report.download');

    Route::get('/email/index', [MailController::class, 'index'])->name('mail.index');
    Route::post('/email/send', [MailController::class, 'send'])->name('mail.send');
});
Route::get('/metrics', function () {
    $adapter = new Redis([
        'host' => env('REDIS_HOST', 'redis'),
        'port' => env('REDIS_PORT', 6379),
    ]);
    $registry = new CollectorRegistry($adapter);

    $renderer = new RenderTextFormat();

    try {
        $metrics = $renderer->render($registry->getMetricFamilySamples());
    } catch (\Exception $e) {
        return response("Error rendering metrics: ".$e->getMessage(), 500);
    }

    return response($metrics, 200)
        ->header('Content-Type', RenderTextFormat::MIME_TYPE);
});
