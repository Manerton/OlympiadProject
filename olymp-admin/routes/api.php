<?php

use App\Http\Controllers\api\ApplicationApiController;
use App\Http\Controllers\api\EventApiController;
use App\Http\Controllers\api\ParticipantApiController;
use App\Http\Controllers\api\ReportApiController;
use App\Http\Controllers\api\SchoolApiController;
use App\Http\Controllers\api\UserApiController;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Route;
use Prometheus\CollectorRegistry;
use Prometheus\RenderTextFormat;
use Prometheus\Storage\Redis;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
|
| Here is where you can register API routes for your application. These
| routes are loaded by the RouteServiceProvider within a group which
| is assigned the "api" middleware group. Enjoy building your API!
|
*/

Route::middleware('auth:sanctum')->get('/user', function (Request $request) {
    return $request->user();
});
Route::group(['middleware' => 'prometheus'], function() {
    Route::get('/user/index/{page?}', [UserApiController::class, 'index'])->name('user-api.index');
    Route::get('/user/create', [UserApiController::class, 'create'])->name('user-api.create');
    Route::post('/user/store', [UserApiController::class, 'store'])->name('user-api.store');
    Route::get('/user/edit/{id}', [UserApiController::class, 'edit'])->name('user-api.edit');
    Route::put('/user/update/{id}', [UserApiController::class, 'update'])->name('user-api.update');
    Route::delete('/user/delete/{id}', [UserApiController::class, 'delete'])->name('user-api.delete');
    Route::get('/user/show/{id}', [UserApiController::class, 'show'])->name('user-api.show');
    Route::post('/user/revoke/{id}', [UserApiController::class, 'revoke'])->name('user.revoke');

    Route::get('/participant/index/{page?}', [ParticipantApiController::class, 'index'])->name('participant-api.index');
    Route::get('/participant/create', [ParticipantApiController::class, 'create'])->name('participant-api.create');
    Route::post('/participant/store', [ParticipantApiController::class, 'store'])->name('participant-api.store');
    Route::get('/participant/edit/{id}', [ParticipantApiController::class, 'edit'])->name('participant-api.edit');
    Route::post('/participant/update/{id}', [ParticipantApiController::class, 'update'])->name('participant-api.update');
    Route::delete('/participant/delete/{id}', [ParticipantApiController::class, 'delete'])->name('participant-api.delete');
    Route::get('/participant/show/{id}', [ParticipantApiController::class, 'show'])->name('participant-api.show');

    Route::get('/school/index/{page?}', [SchoolApiController::class, 'index'])->name('school-api.index');
    Route::get('/school/create', [SchoolApiController::class, 'create'])->name('school-api.create');
    Route::post('/school/store', [SchoolApiController::class, 'store'])->name('school-api.store');
    Route::get('/school/edit/{id}', [SchoolApiController::class, 'edit'])->name('school-api.edit');
    Route::post('/school/update/{id}', [SchoolApiController::class, 'update'])->name('school-api.update');
    Route::delete('/school/delete/{id}', [SchoolApiController::class, 'delete'])->name('school-api.delete');
    Route::get('/school/show/{id}', [SchoolApiController::class, 'show'])->name('school-api.show');

    Route::get('/report/index', [ReportApiController::class, 'index'])->name('report-api.index');
    Route::get('/report/download/{id}', [ReportApiController::class, 'download'])->name('report-api.download');

    Route::get('/application/index/{page?}', [ApplicationApiController::class, 'index'])->name('application-api.index');
    Route::get('/application/create', [ApplicationApiController::class, 'create'])->name('application-api.create');
    Route::post('/application/store', [ApplicationApiController::class, 'store'])->name('application-api.store');
    Route::get('/application/edit/{id}', [ApplicationApiController::class, 'edit'])->name('application-api.edit');
    Route::put('/application/update/{id}', [ApplicationApiController::class, 'update'])->name('application-api.update');
    Route::delete('/application/delete/{id}', [ApplicationApiController::class, 'delete'])->name('application-api.delete');
    Route::get('/application/show/{id}', [ApplicationApiController::class, 'show'])->name('application-api.show');
    Route::post('/application/confirm/{id}', [ApplicationApiController::class, 'confirm'])->name('application-api.confirm');
    Route::post('/application/reject/{id}', [ApplicationApiController::class, 'reject'])->name('application-api.reject');

    Route::get('/event/index/{page?}', [EventApiController::class, 'index'])->name('event-api.index');
    Route::get('/event/show/{id}', [EventApiController::class, 'show'])->name('event-api.show');
    Route::delete('/event/delete/{id}', [EventApiController::class, 'delete'])->name('event-api.delete');
    Route::get('/event/task/{id}', [EventApiController::class, 'task'])->name('event-api.task');
    Route::get('/event/attendance/{id}', [EventApiController::class, 'attendance'])->name('event-api.attendance');
    Route::get('/event/point/{id}', [EventApiController::class, 'point'])->name('event-api.point');
    Route::get('/event/synchronize/{id}', [EventApiController::class, 'synchronize'])->name('event-api.synchronize');
    Route::post('/event/add-task/{id}', [EventApiController::class, 'addTask'])->name('event-api.add-task');
    Route::post('/event/change-attendance', [EventApiController::class, 'changeAttendance'])->name('event-api.change-attendance');
    Route::delete('/event/delete-task/{id}', [EventApiController::class, 'deleteTask'])->name('event-api.delete-task');
    Route::post('/event/change-score', [EventApiController::class, 'changeScore'])->name('event-api.change-score');
    Route::get('/event/prize-score/{id}', [EventApiController::class, 'prizeScore'])->name('event-api.prize-score');
    Route::post('/event/set-prize-score/{id}', [EventApiController::class, 'setPrizeScore'])->name('event-api.set-prize-score');
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
